from fastapi import FastAPI


def test_main_exposes_fastapi_app() -> None:
    from app.main import app

    assert isinstance(app, FastAPI)
