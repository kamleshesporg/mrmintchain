import pytest

from .network import setup_mrmintchain, setup_geth


@pytest.fixture(scope="session")
def mrmintchain(tmp_path_factory):
    path = tmp_path_factory.mktemp("mrmintchain")
    yield from setup_mrmintchain(path, 26650)


@pytest.fixture(scope="session")
def geth(tmp_path_factory):
    path = tmp_path_factory.mktemp("geth")
    yield from setup_geth(path, 8545)


@pytest.fixture(
    scope="session", params=["mrmintchain", "mrmintchain-ws"]
)
def mrmintchain_rpc_ws(request, mrmintchain):
    """
    run on both mrmintchain and mrmintchain websocket
    """
    provider = request.param
    if provider == "mrmintchain":
        yield mrmintchain
    elif provider == "mrmintchain-ws":
        mrmintchain_ws = mrmintchain.copy()
        mrmintchain_ws.use_websocket()
        yield mrmintchain_ws
    else:
        raise NotImplementedError
