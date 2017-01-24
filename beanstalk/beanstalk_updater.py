# -*- coding: utf-8 -*-

from tracer import Tracer
import boto3
import botocore
import echoes
import os

echo = echoes.get_echo()

@Tracer()
def update_beanstalk(zip_file, configs, user_options):

    @Tracer()
    def send_zip_file_to_s3(zip_file):
        """
        Send the given zip_file to the bucket set on the configuration file
        """
        basename = os.path.join(
            configs['projects'][user_options.project]['environments'][user_options.environment]['s3_directory'],
            os.path.basename(zip_file))
        bucket = configs['projects'][user_options.project]['environments'][user_options.environment]['s3_bucket']

        echo.info('Will send %s file to %s bucket at %s location' %
                    (zip_file, bucket, basename))

        if not user_options.dryrun:
            s3 = boto3.resource('s3')
            s3.Object(bucket, basename).put(Body=open(zip_file, 'rb'))

        echo.info('File sent to S3')
        echo.debug('File sent to S3 on bucket {} with path {}'.format(bucket, basename))

        return bucket, basename

    @Tracer()
    def register_newest_zip_as_application_version(bucket, key):
        sha = configs['projects'][user_options.project]['sha']
        app_name = configs['projects'][user_options.project]['environments'][user_options.environment]['beanstalk-app']
        version_label = configs['projects'][user_options.project]['sha'] + '-' + user_options.environment + '-' + user_options.project
        description = 'Version {} of {} built automatically'.format(sha, user_options.project)
        source_bundle = {'S3Bucket': bucket, 'S3Key': key}

        @Tracer()
        def delete_application_version(client, app_name, version_label):
            if not user_options.dryrun:
                client.delete_application_version(ApplicationName=app_name, VersionLabel=version_label)

        @Tracer()
        def create_application_version(client, app_name, version_label, description, source_bundle):
            if not user_options.dryrun:
                client.create_application_version(ApplicationName=app_name,
                                                  VersionLabel=version_label,
                                                  Description=description,
                                                  SourceBundle=source_bundle)

        beanstalk_client = boto3.client('elasticbeanstalk', region_name=configs['region'])
        try:
            try:
                echo.debug('Creating version {} to send to Beanstalk as an application version'.format(sha))
                create_application_version(beanstalk_client, app_name, version_label, description, source_bundle)
            except botocore.exceptions.ClientError as e:
                echo.debug('Version {} alread exists as an application version'.format(sha))
                delete_application_version(beanstalk_client, app_name, version_label)
                echo.debug('Creating version {} again to send to Beanstalk as an application version'.format(sha))
                create_application_version(beanstalk_client, app_name, version_label, description, source_bundle)
        except Exception as e:
            echo.warn('An error occurred while trying to register the newest zip as an application version')
            raise e

        echo.info('Version {} sent to Beanstalk as an application version'.format(sha))

        return version_label

    @Tracer()
    def update_environment(version_label):
        if user_options.deploy and not user_options.dryrun:
            echo.warn('Will deploy newest version to ElasticBeanstalk environment.')

            app_name = configs['projects'][user_options.project]['environments'][user_options.environment]['beanstalk-app']
            env_name = configs['projects'][user_options.project]['environments'][user_options.environment]['beanstalk-env']

            beanstalk_client = boto3.client('elasticbeanstalk', region_name=configs['region'])
            echo.info('Issuing the environment update to beanstalk')
            try:
                response = beanstalk_client.update_environment(ApplicationName=app_name, EnvironmentName=env_name, VersionLabel=version_label)
                echo.debug('Response from the update_environment call:')
                echo.debug(response)
                echo.info('>>>')
                echo.info('Environment Update issued. Please check at Elastic Beanstalk page to see the progress of the update.')
                echo.info('If something goes wrong and you want a fast retry, add the --cache to the previous command line.')
                echo.info('>>>')
            except Exception as e:
                echo.warn('An error has occurred when issuing the environment update on ElasticBeanstalk:')
                raise e

        else:
            echo.info('Will not deploy the application to the ElasticBeanstalk environment.')

    bucket, key = send_zip_file_to_s3(zip_file)
    version_label = register_newest_zip_as_application_version(bucket, key)
    update_environment(version_label)
