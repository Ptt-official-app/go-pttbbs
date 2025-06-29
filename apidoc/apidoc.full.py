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


@app.route(_with_app_prefix('/class/<cls>/board'), methods=['POST'])
def _create_board(cls):
    """
    swagger_from_file: apidoc/create_board.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/autocomplete'), methods=['GET'])
def _load_auto_complete_boards():
    """
    swagger_from_file: apidoc/load_auto_complete_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards'), methods=['GET'])
def _load_general_boards():
    """
    swagger_from_file: apidoc/load_general_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/detail'), methods=['GET'])
def _load_general_board_details():
    """
    swagger_from_file: apidoc/load_general_board_details.yaml
    """
    return ''


@app.route(_with_app_prefix('/cls/<clsid>/boards'), methods=['GET'])
def _load_class_boards(clsid):
    """
    swagger_from_file: apidoc/load_class_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/cls/boards'), methods=['GET'])
def _load_full_class_boards():
    """
    swagger_from_file: apidoc/load_full_class_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/articles'), methods=['GET'])
def _load_general_articles(bid):
    """
    swagger_from_file: apidoc/load_general_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/detail'), methods=['GET'])
def _get_board_detail(bid):
    """
    swagger_from_file: apidoc/get_board_detail.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/articles/bottom'), methods=['GET'])
def _load_bottom_articles(bid):
    """
    swagger_from_file: apidoc/load_bottom_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/summary'), methods=['GET'])
def _load_board_summary(bid):
    """
    swagger_from_file: apidoc/load_board_summary.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/popular'), methods=['GET'])
def _load_hot_boards():
    """
    swagger_from_file: apidoc/load_hot_boards.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/byclass'), methods=['GET'])
def _load_general_boards_by_class():
    """
    swagger_from_file: apidoc/load_general_boards_by_class.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/bybids'), methods=['POST'])
def _load_boards_by_bids():
    """
    swagger_from_file: apidoc/load_boards_by_bids.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/isvalid'), methods=['GET'])
def _is_board_valid_user(bid):
    """
    swagger_from_file: apidoc/is_board_valid_user.yaml
    """
    return ''


@app.route(_with_app_prefix('/boards/isvalid'), methods=['POST'])
def _is_boards_valid_user():
    """
    swagger_from_file: apidoc/is_boards_valid_user.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/posttemplate/<tid>'), methods=['GET'])
def _get_post_template(bid, tid):
    """
    swagger_from_file: apidoc/get_post_template.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article/<aid>'), methods=['GET'])
def _get_article(bid, aid):
    """
    swagger_from_file: apidoc/get_article.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article'), methods=['POST'])
def _create_article(bid):
    """
    swagger_from_file: apidoc/create_article.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article/<aid>/edit'), methods=['POST'])
def _edit_article(bid, aid):
    """
    swagger_from_file: apidoc/edit_article.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/deletearticles'), methods=['POST'])
def _delete_articles(bid, aid):
    """
    swagger_from_file: apidoc/delete_articles.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article/<aid>/comment'), methods=['POST'])
def _create_comment(bid, aid):
    """
    swagger_from_file: apidoc/create_comment.yaml
    """
    return ''


@app.route(_with_app_prefix('/board/<bid>/article/<aid>/crosspost'), methods=['POST'])
def _cross_post(bid, aid):
    """
    swagger_from_file: apidoc/cross_post.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/information'), methods=['GET'])
def _get_user(uid):
    """
    swagger_from_file: apidoc/get_user.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/changepasswd'), methods=['POST'])
def _change_passwd(uid):
    """
    swagger_from_file: apidoc/change_passwd.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/changeemail'), methods=['POST'])
def _change_email(uid):
    """
    swagger_from_file: apidoc/change_email.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/attemptchangeemail'), methods=['POST'])
def _attempt_change_email(uid):
    """
    swagger_from_file: apidoc/attempt_change_email.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/attemptsetidemail'), methods=['POST'])
def _attempt_set_id_email(uid):
    """
    swagger_from_file: apidoc/attempt_set_id_email.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/setidemail'), methods=['POST'])
def _set_id_email(uid):
    """
    swagger_from_file: apidoc/set_id_email.yaml
    """
    return ''


@app.route(_with_app_prefix('/token/info'), methods=['POST'])
def _get_token_info():
    """
    swagger_from_file: apidoc/get_token_info.yaml
    """
    return ''


@app.route(_with_app_prefix('/emailtoken/info'), methods=['POST'])
def _get_email_token_info():
    """
    swagger_from_file: apidoc/get_email_token_info.yaml
    """
    return ''


@app.route(_with_app_prefix('/refresh'), methods=['POST'])
def _refresh():
    """
    swagger_from_file: apidoc/refresh.yaml
    """
    return ''


@app.route(_with_app_prefix('/refreshtoken/info'), methods=['POST'])
def _get_refresh_token_info():
    """
    swagger_from_file: apidoc/get_refresh_token_info.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/favorites'), methods=['GET'])
def _get_fav(uid):
    """
    swagger_from_file: apidoc/get_fav.yaml
    """
    return ''


@app.route(_with_app_prefix('/user/<uid>/favorites/post'), methods=['POST'])
def _write_favorites(uid):
    """
    swagger_from_file: apidoc/write_favorites.yaml
    """
    return ''


@app.route(_with_app_prefix('/uservisitcount'), methods=['GET'])
def _get_user_visit_count():
    """
    swagger_from_file: apidoc/get_user_visit_count.yaml
    """
    return ''


@app.route(_with_app_prefix('/existsuser'), methods=['POST'])
def _check_exists_user():
    """
    swagger_from_file: apidoc/check_exists_user.yaml
    """
    return ''


@app.route(_with_app_prefix('/version'), methods=['GET'])
def _version():
    """
    swagger_from_file: apidoc/version.yaml
    """
    return ''


@app.route(_with_app_prefix('/admin/reloaduhash'), methods=['GET'])
def _reloaduhash():
    """
    swagger_from_file: apidoc/reload_uhash.yaml
    """
    return ''


@app.route(_with_app_prefix('/admin/user/<uid>/setperm'), methods=['POST'])
def _set_user_perm(uid):
    """
    swagger_from_file: apidoc/set_user_perm.yaml
    """
    return ''
