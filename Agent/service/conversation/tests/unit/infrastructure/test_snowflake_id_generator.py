from app.infrastructure.id_generator.snowflake import SnowflakeIdGenerator


def test_snowflake_id_generator_returns_monotonic_signed_int64_ids() -> None:
    now = 1_800_000_000_000
    generator = SnowflakeIdGenerator(
        node_id=7,
        epoch_ms=1_735_689_600_000,
        time_ms=lambda: now,
    )

    first_id = generator.next_id()
    second_id = generator.next_id()

    assert isinstance(first_id, int)
    assert first_id < second_id
    assert 0 <= first_id < 2**63
    assert 0 <= second_id < 2**63
