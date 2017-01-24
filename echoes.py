# -*- coding: utf-8 -*-
import click
from pprint import pformat

_echo_object = None

class Echo():
    def __init__(self, debug_mode=False):
        """
        :debug_mode: Boolean
        """
        self.debug_mode = debug_mode

    def _echo(self, msg, color, kind):
        message = msg

        if type(msg) is dict:
            message = pformat(message)
        message = message.strip('\n')

        click.echo(click.style(kind + ': ' + message, fg=color))

    def info(self, msg):
        self._echo(msg, 'white', 'INFO')

    def warn(self, msg):
        self._echo(msg, 'yellow', 'WARN')

    def error(self, msg):
        self._echo(msg, 'red', 'ERROR')

    def debug(self, msg):
        if self.debug_mode:
            self._echo(msg, 'cyan', 'DEBUG')

def get_echo(debug_mode=False):
    """
    It will return the same Echo object independently from where it was called
    and create a new one if there is none. This is very useful as one can
    define the debug_mode of the entire program in a single place an then use
    it everywhere.

    Usage:

    # On the main file:

    import echoes
    echo = echoes.get_echo()     # To dont't display debug messages
    # OR
    echo = echoes.get_echo(True) # To display debug messages


    # On all the other files:

    import echoes
    echo = echoes.get_echo()
    """
    global _echo_object
    if not _echo_object:
        _echo_object = Echo(debug_mode)
    return _echo_object
