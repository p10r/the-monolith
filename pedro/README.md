# Pedro

A little tool to notify you via Telegram that your beloved artists are playing a gig in your
area.

## TODOs:

Product:
- Sort events by artist
- Filter out duplicate events (e.g. to promoters added the same event) 

Tech:
- parallel tests
- short/long tests
- handle 404 when requesting events
- throw error if artist cant be found
- indicate if there's a space
- Use JSON functionality of sqlite
- Give user info if they don't follow anyone yet
- Create RAError and TelegramError to improve logs