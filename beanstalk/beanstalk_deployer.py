# -*- coding: utf-8 -*-

from distutils import dir_util
from tracer import Tracer
from . import beanstalk_file_creator
from . import beanstalk_updater
import distutils
import echoes
import os
import tempfile

echo = echoes.get_echo()

@Tracer()
def deploy_to_beanstalk(project_directory, configs, user_options):
    """
    Deploy the newly built image to ElasticBeanstalk's Environment
    """

    # INFO: Interesting repo that we can get some refs and ideas: https://github.com/briandilley/ebs-deploy
    temp_dir_path = os.path.join(tempfile.gettempdir(), user_options.environment + '-' + user_options.project + '-' + configs['projects'][user_options.project]['sha'])

    if os.path.exists(temp_dir_path):
        echo.info('Directory already exists: %s ...cleaning it.' % str(temp_dir_path))
        distutils.dir_util.remove_tree(temp_dir_path)
        echo.info('Directory removed.')

    os.makedirs(temp_dir_path)
    echo.debug('Temporary directory created: %s' % str(temp_dir_path))

    zip_file = beanstalk_file_creator.create_zip_file(project_directory, temp_dir_path, configs, user_options)
    beanstalk_updater.update_beanstalk(zip_file, configs, user_options)
