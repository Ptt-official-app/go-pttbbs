# -*- coding: utf-8 -*-

import yaml
from flask import Flask

app = Flask(__name__.split('.')[0])

_APP_PREFIX = '/v1'


def _with_app_prefix(path):
    """Summary

    Args:
        path (TYPE): Description

    Returns:
        TYPE: Description
    """
    if not _APP_PREFIX:
        return path
    return _APP_PREFIX + path


with open('apidoc/template.json', 'r') as f:
    template = yaml.load(f)


@app.route('/', methods=['POST'])
def _index():
    """
    swagger_from_file: apidoc/index.yaml
    """
    return ''


@app.route(_with_app_prefix('/users/me/information'), methods=['POST'])
def _me():
    """
    swagger_from_file: apidoc/me.yaml
    """
    return ''


@app.route(_with_app_prefix('/token'), methods=['POST'])
def _login():
    """
    swagger_from_file: apidoc/login.yaml
    """
    return ''


@app.route(_with_app_prefix('/register'), methods=['POST'])
def _register():
    """
    swagger_from_file: apidoc/register.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/boards'), methods=['POST'])
def _load_general_boards():
    """
    swagger_from_file: apidoc/load_general_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/articles'), methods=['POST'])
def _load_general_articles(bid):
    """
    swagger_from_file: apidoc/load_general_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/class/<class_id>'), methods=['POST'])
def _load_boards_by_class(class_id):
    """
    swagger_from_file: apidoc/load_boards_by_class.yaml
    """
    return ''
