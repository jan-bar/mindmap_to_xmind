#!/usr/bin/env python
#-*- coding: utf-8 -*-

"""
    xmind.core.title
    ~~~~~~~~~~~~~~~

    :copyright:
    :license:

"""

__author__ = "aiqi@xmind.net <Woody Ai>"

from . import const

from .mixin import WorkbookMixinElement


class TitleElement(WorkbookMixinElement):
    TAG_NAME = const.TAG_TITLE

    def __init__(self, node, ownerWorkbook):
        super(TitleElement, self).__init__(node, ownerWorkbook)

