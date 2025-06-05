from flask import Blueprint, jsonify

bp = Blueprint('routes', __name__)

@bp.route('/status')
def status():
    return jsonify({"status": "API running"})
