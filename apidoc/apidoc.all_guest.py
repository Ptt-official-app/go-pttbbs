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
    template = yaml.full_load(f)


@app.route('/', methods=['POST'])
def _index():
    """
    swagger_from_file: apidoc/index.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/summary'), methods=['GET'])
def _load_board_summary(bid):
    """
    swagger_from_file: apidoc/load_board_summary.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/articles'), methods=['GET'])
def _load_general_articles(bid):
    """
    swagger_from_file: apidoc/load_general_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/articles/bottom'), methods=['GET'])
def _load_bottom_articles(bid):
    """
    swagger_from_file: apidoc/load_bottom_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article/<aid>'), methods=['GET'])
def _get_article(bid, aid):
    """
    swagger_from_file: apidoc/get_article.yaml
    """
    return ''


@app.route(_with_app_prefix('/version'), methods=['GET'])
def _version():
    """
    swagger_from_file: apidoc/version.yaml
    """
    return ''
