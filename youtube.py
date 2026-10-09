import aiohttp, asyncio, psycopg, json, urllib.parse, db, utils, youtube_api


async def get_new_channels_metadata():
	async with await psycopg.AsyncConnection.connect("dbname=media_monitoring", row_factory = psycopg.rows.dict_row, autocommit = True) as db_connection:

		channels = await db.get_channels(db_connection)
		updates = []

		async with aiohttp.ClientSession(timeout = aiohttp.ClientTimeout(total = 10), headers = {"Accept-Encoding": "gzip"}) as session:
			for start in range(0, len(channels), 1000):
				batch = channels[start:start + 1000]

				results = await asyncio.gather(*(youtube_api.get_channel_metadata(channel["channel_id"], session) for channel in batch), return_exceptions = True)

				for channel, metadata in zip(batch, results):
					if isinstance(metadata, Exception):
						metadata = {"metadata_error": f"{type(metadata).__name__}: {metadata}"}
					
					updates.append({**channel, **{key: value for key, value in metadata.items() if value is not None}, "metadata_error": metadata["metadata_error"]})

				print(f"Processed: {start + len(batch)}/{len(channels)}")

			await db.update_channels_metadata(db_connection, updates)


asyncio.run(get_new_channels_metadata())