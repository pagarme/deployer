# -*- coding: utf-8 -*-

from git import Repo, Remote
from tracer import Tracer
import os
import echoes
echo = echoes.get_echo()

@Tracer()
def clone_repositories(root_dir, project, configs, user_options):
    """
    Clone the necessary repositories for our process
    """

    @Tracer()
    def find_and_set_sha(repo):
        """
        Gets the SHA of a repo, 'latest' if there is no repo (usually at a
        dryrun) and sets it at the config object for later use.

        :repo: The repository to find the sha to use
        :returns: The SHA that will be used as tag on the built docker image
        """
        latest_sha = 'latest'
        if repo:
            latest_sha = repo.git.rev_parse('HEAD', short=7)
        echo.debug('Repo not set, will use %s as default SHA' % latest_sha)

        configs['projects'][user_options.project]['sha'] = latest_sha  # Not a good practice I know :/
        echo.debug('Latest SHA: %s' % latest_sha)


    echo.info('Will clone the repositories')

    repo_name = project['repo']
    repos = configs['repositories']
    repo = repos[repo_name]

    git_url = repo['git_url']
    branch = repo['branch']
    repo_dir = os.path.join(os.path.abspath(root_dir), repo_name)
    repo['repo_dir'] = repo_dir  # Not a good practice I know :/

    repo = None
    if not user_options.dryrun:
        if os.path.exists(repo_dir):
            # INFO: I'm still not too happy with this solution, there
            # is still the edge case that if someone pushed new things
            # it may break because we haven't merged.
            echo.info('Fetching %s ...' % git_url)

            echo.debug('Fetching origin from {} into {}'.format(git_url, repo_dir))
            repo = Repo(repo_dir)
            Remote(repo, 'origin').fetch()

            echo.debug('Checking out branch {}'.format(branch))
            git = repo.git
            git.checkout('--force', branch)
        else:
            echo.info('Cloning {} into {}'.format(git_url, repo_dir))
            repo = Repo.clone_from(git_url, repo_dir,
                                    recursive=True, branch=branch)

    find_and_set_sha(repo)
