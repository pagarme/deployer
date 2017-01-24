# -*- coding: utf-8 -*-

from subprocess import Popen, PIPE
from tracer import Tracer
import echoes
import os
import shlex

echo = echoes.get_echo()


@Tracer()
def build_images(applications_directory, project, configs, user_options):
    built_images = []
    images_to_build_paths = []

    if project['multiproject']:
        for inner_project in project['inner_projects']:
            project_to_build = None

            if isinstance(inner_project, str):
                project_to_build = inner_project
            else:
                project_name = list(inner_project)[0]
                project_to_build = list(inner_project.values())[0]['main_project']

            images_to_build_paths.append(os.path.join(applications_directory, project_to_build))
    else:
        images_to_build_paths.append(os.path.join(applications_directory, user_options.project))

    for image_path in images_to_build_paths:
        if not image_path in built_images:
            build_image(image_path, project, configs, user_options)
            echo.debug('Adding {} to the list of built images'.format(image_path))
            built_images.append(image_path)
        else:
            echo.info('Will not build image looking at "{}" because it was already built'.format(image_path))

@Tracer()
def build_image(work_dir, project, configs, user_options):
    """
    Build the docker image of the application
    """
    echo.info('Will build image for "{}" and environment "{}", looking for files in "{}"'.format(user_options.project, user_options.environment, work_dir))

    repo_name = project['repo']
    repo = configs['repositories'][repo_name]
    repository_path = repo['repo_dir']

    cache_param = '' if user_options.cache else '--no-cache'
    push_param = '--push' if user_options.push else ''

    build_vars = {}
    build_vars['ENVIRONMENT'] = user_options.environment
    build_vars['REPOSITORY_PATH'] = repository_path
    build_vars['IMAGE_SHA'] = configs['projects'][user_options.project]['sha']

    vars_param = ''
    for key, value in build_vars.items():
        vars_param += '--var "{}={}" '.format(key, value)

    command = 'rocker build {} {} {}'.format(cache_param, push_param, vars_param).strip()

    echo.debug("Will execute command '{}'\nat the work dir '{}'\n".format(command, work_dir))

    if not user_options.dryrun:
        err = None

        try:
            process = Popen(shlex.split(command),
                            cwd=work_dir,
                            stdout=PIPE,
                            stderr=PIPE)

            echo.info('Command stdout:')
            for line in iter(process.stdout.readline, ''):
                echo.info(str(line, 'utf-8'))
                if not line:
                    break

            _, err = process.communicate()
            return_code = process.returncode

            if return_code != 0:
                raise Exception('Rocker returned a status code different from 0. Please check.')
        except Exception as e:
            echo.warn('An error occurred while trying to build the image!')
            raise e
        finally:
            echo.debug('Last command execution stderr:')
            echo.debug(str(err))
            echo.debug('')

