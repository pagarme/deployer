# -*- coding: utf-8 -*-

import echoes
echo = echoes.get_echo()

class Tracer(object):
    def __call__(self, func):
        from itertools import chain
        def wrapper(that, *args, **kwargs):
            name = func.__name__
            echo.debug("TRACER - Calling: {}({})".format(name, ", ".join(map(repr, chain(args, kwargs.values())))))
            return func(that, *args, **kwargs)
        return wrapper
