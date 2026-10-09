

async def get_channels(db_connection):
	return await (await db_connection.execute("SELECT * FROM youtube_channels WHERE metadata_updated_at IS NULL AND metadata_error IS NULL")).fetchall()


async def update_channels_metadata(db_connection, channels):
	async with db_connection.transaction():
		async with db_connection.cursor() as cursor:
			await cursor.executemany("UPDATE youtube_channels SET name = %(name)s, description = %(description)s, username = %(username)s, avatar_url = %(avatar_url)s, subscriber_count = %(subscriber_count)s, video_count = %(video_count)s, metadata_updated_at = CASE WHEN %(metadata_error)s::text IS NULL THEN NOW() ELSE metadata_updated_at END, metadata_error = %(metadata_error)s WHERE channel_id = %(channel_id)s", channels)