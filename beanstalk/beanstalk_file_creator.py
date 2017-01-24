# -*- coding: utf-8 -*-

from tracer import Tracer
import distutils
import echoes
import os
import zipfile

echo = echoes.get_echo()

SEPARATOR = os.sep  # /

@Tracer()
def build_context(configs, user_options):
    context = {
        'LOGENTRIES_TOKEN': configs['projects'][user_options.project]['environments'][user_options.environment]['logentries_token']
    }

    project = configs['projects'][user_options.project]
    if project['multiproject']:
        for inner_project in project['inner_projects']:
            if isinstance(inner_project, str):
                # input (part of the configs.yml file):
                #
                # inner_projects:
                # - hookshot
                #
                # output:
                #
                # hookshot
                project_name = inner_project
            else:
                # input (part of the configs.yml file):
                #
                # inner_projects:
                # - hookshot:
                #   main_project: 'api'
                #
                # output:
                #
                # hookshot
                project_name = list(inner_project)[0]

            context[project_name.upper() + "_IMAGE_SHA"] = configs['projects'][user_options.project]['sha']
            context[project_name.upper() + '_PAGARME_ENV'] = user_options.environment

    else:
        context['IMAGE_SHA'] = configs['projects'][user_options.project]['sha']
        context['PAGARME_ENV'] = user_options.environment

    return context


@Tracer()
def create_zip_file(project_directory, temp_dir_path, configs, user_options):

    @Tracer()
    def copy_ebextensions(temp_path):
        """
        Copies the .ebextensions folder to the temp_path. So it later on
        will be zipped on the file that will be sent to S3.
        """
        src = os.path.join(project_directory, '.ebextensions')
        if os.path.exists(src):
            dst = os.path.join(temp_path, '.ebextensions')
            distutils.dir_util.copy_tree(src, dst)
        else:
            echo.warn('No .ebextensions directory to copy.')
            echo.debug('Looked at path: %s' % str(project_directory))

    @Tracer()
    def modify_dockerrun(path):
        """
        Modifies the Dockerrun.aws.json file to use the newly built image.
        It changes the '{{ IMAGE_SHA }}' that is defined in the file.

        :path: Where the file will be saved
        """
        from jinja2 import Environment, Template, FileSystemLoader
        template_loader = FileSystemLoader(SEPARATOR)
        template_env = Environment(loader=template_loader)
        original_file = os.path.join(project_directory, 'Dockerrun.aws.json')
        original_file = SEPARATOR.join(original_file.split(SEPARATOR)[1:])

        echo.debug('Will get the Dockerrun file from: %s' %
                    original_file)
        template = template_env.get_template(original_file)

        context = build_context(configs, user_options)

        echo.debug('Values that will be changed on the template: %s' % context)

        dockerrun_path = os.path.join(path, 'Dockerrun.aws.json')
        with open(dockerrun_path, 'w') as f:
            f.write(template.render(context))

        echo.debug('Modified Dockerrun.aws.json in: %s' % dockerrun_path)

    @Tracer()
    def create_archive(directory, filename, ignored_files=['.git', '.svn']):
        """
        Creates a zip file from the given directory with the given filename as name.
        """
        with zipfile.ZipFile(filename, 'w', compression=zipfile.ZIP_DEFLATED) as zip_file:
            root_len = len(os.path.abspath(directory))

            # create it
            echo.info("Creating archive: " + str(filename))
            for root, dirs, files in os.walk(directory, followlinks=True):
                archive_root = os.path.abspath(root)[root_len + 1:]
                for f in files:
                    fullpath = os.path.join(root, f)
                    archive_name = os.path.join(archive_root, f)

                    # ignore the file we're creating
                    if filename in fullpath:
                        continue

                    # ignored files
                    if ignored_files is not None:
                        for name in ignored_files:
                            if fullpath.endswith(name):
                                echo.info("Skipping: " + str(name))
                                continue

                    echo.info("Adding: " + str(archive_name))
                    zip_file.write(fullpath, archive_name, zipfile.ZIP_DEFLATED)

        return filename

    copy_ebextensions(temp_dir_path)
    modify_dockerrun(temp_dir_path)
    return create_archive(temp_dir_path, temp_dir_path + '.zip')
