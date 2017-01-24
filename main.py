#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from user_options import UserOptions
import click
import echoes
import os
import sys
import tempfile
import traceback
import yaml

CONFIGS_PATH = 'configs.yml'
SEPARATOR = os.sep  # /
UPDIR = os.pardir  # ..

@click.command()
@click.option('--cache/--no-cache',
              default=False,
              help='If rocker should use cache for its builds or not. It doesnt use cache by default.')
@click.option('--debug',
              is_flag=True,
              help='If you want to activate debug messages.')
@click.option('--deploy/--no-deploy',
              default=False,
              help='If you really want to deploy this new release to ElasticBeanstalk. Default = False')
@click.option('--dryrun', '-d',
              is_flag=True)
@click.option('--environment', '-e',
              required=True,
              default='main')
@click.option('--project', '-p',
              required=True)
@click.option('--push/--no-push',
              default=True,
              help='If the docker images should be pushed. Default = True')
def main(cache, debug, deploy, dryrun, environment, project, push):
    echo = echoes.get_echo(debug)

    def read_configs():
        """
        Read configuration from the configuration file
        """
        echo.info('Reading configuration file...')
        try:
            with open(CONFIGS_PATH, 'r') as yml_file:
                configs = yaml.load(yml_file)

            return configs
        except Exception as e:
            msg = 'An error ocurred while trying to read configuration file {}:\n'.format(
                CONFIGS_PATH)
            msg += str(e)

            echo.warn(msg)
            sys.exit(1)


    configs = read_configs()

    user_options = UserOptions(
        cache=cache,
        debug=debug,
        deploy=deploy,
        dryrun=dryrun,
        environment=environment,
        project=project,
        push=push
    )

    projects = configs['projects']

    echo.debug('DEBUG MODE ACTIVATED')

    # Validate if options are valid
    if not user_options.project in projects:
        echo.error('Project {} not found'.format(user_options.project))
        sys.exit(1)

    project = projects[user_options.project]
    environments = project['environments']
    
    if not user_options.environment in environments:
        echo.error('Environment {} not found'.format(user_options.project))
        sys.exit(1)

    environment = environments[user_options.environment]

    try:
        root_dir = os.path.join(tempfile.gettempdir(), 'deployer_workspace')
        if not os.path.exists(root_dir):
            os.makedirs(root_dir)

        file_dir = os.path.dirname(__file__)
        applications_directory = os.path.abspath(os.path.join(file_dir, UPDIR, 'applications'))
        project_directory = os.path.join(applications_directory, user_options.project)

        import cloner
        cloner.clone_repositories(root_dir,
                                  project,
                                  configs,
                                  user_options)

        import image_builder
        image_builder.build_images(applications_directory,
                                   project,
                                   configs,
                                   user_options)
        sys.exit(1)

        from beanstalk import beanstalk_deployer
        beanstalk_deployer.deploy_to_beanstalk(project_directory,
                                               configs,
                                               user_options)
    except Exception as e:
        echo.warn(
            'An error has occurred while trying to deploy the newest version of our application')
        echo.warn(str(e))

        for frame in traceback.extract_tb(sys.exc_info()[2]):
            fname,lineno,fn,text = frame
            echo.warn("Error in {} on line {}".format(fname, lineno))

        sys.exit(1)


if __name__ == "__main__":
    main()
