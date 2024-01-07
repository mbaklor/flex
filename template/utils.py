import logging
from logging.handlers import RotatingFileHandler
from platform import python_version
from json import load
from uci import Uci  # type: ignore

M400 = "M400"
MPI400 = "MPI400"
LX400 = "LX400"


# implement a rotating log with n log files maximum
def init_logger(
    log_file_name: str, max_file_size: int, num_backups: int, debug_mode: bool
):
    """
    implement a rotating log with n log files maximum
    """
    rot_handler = RotatingFileHandler(
        log_file_name, maxBytes=max_file_size, backupCount=num_backups
    )
    if debug_mode:
        level_mode = logging.DEBUG
    else:
        level_mode = logging.INFO
    logging.basicConfig(
        level=level_mode,
        handlers=[rot_handler],
        format="%(asctime)s\t%(levelname)s\t%(filename)s\t%(message)s",
    )


def get_app_data():
    """
    check runtime info from manifest and uci config

    Returns:
        app_name(str): Name of the app from manifest file
        app_log(str): Name of logfile from manifest, defaults to "app_log.log"
        app_ver(str): Version string from UCI config

    """
    with open("manifest.json", "r") as manifest:
        config = load(manifest)
        app_name: str = config["name"]
        if "app_log" in config:
            app_log: str = config["app_log"]
        else:
            app_log = "app_log.log"
    with Uci() as u:
        app_ver: str = u.get("flexa_agent.service.package_version")
        if not isinstance(app_ver, str):
            app_ver = ""
    return app_name, app_log, app_ver


def get_hw_name():
    """
    check hardware name from factory info

    Returns:
        hw_name(str): Hardware name
    """
    with open("/tmp/factory_info_partition/factory_info_json", "r") as factory_info:
        hw_info = load(factory_info)["HW_DEVICE"]
        hw_name: str = hw_info["Product_Name"].upper()
        if M400 in hw_name:
            return M400
        if MPI400 in hw_name:
            return MPI400
        if LX400 in hw_name:
            return LX400
        return hw_name


def init():
    """
    general init function for flexa:
    starts a logger to the log file, logs out hardware name, python version, app name, app version
    """
    app_name, app_log, app_ver = get_app_data()
    init_logger(f"/var/log/{app_log}", 2 * 1024 * 1024, 3, False)
    logger = logging.getLogger(__name__)
    hw_name = get_hw_name()
    logger.info(
        'Started %s: python=%s app="%s %s"',
        hw_name,
        python_version(),
        app_name,
        app_ver,
    )
