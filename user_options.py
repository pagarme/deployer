# -*- coding: utf-8 -*-

class UserOptions(dict):
    """
    Command line options chosen by the user
    options = UserOptions({'first_name': 'John'}, last_name='Doe', age=42, sports=['Basketball'])
    """

    def __init__(self, *args, **kwargs):
        super(UserOptions, self).__init__(*args, **kwargs)
        for arg in args:
            if isinstance(arg, dict):
                for k, v in arg.iteritems():
                    self[k] = v

        if kwargs:
            for k, v in kwargs.items():
                self[k] = v

    def __getattr__(self, attr):
        return self.get(attr)

    def __setattr__(self, key, value):
        self.__setitem__(key, value)

    def __setitem__(self, key, value):
        super(UserOptions, self).__setitem__(key, value)
        self.__dict__.update({key: value})

    def __delattr__(self, item):
        self.__delitem__(item)

    def __delitem__(self, key):
        super(UserOptions, self).__delitem__(key)
        del self.__dict__[key]
