## go-web-crawler

Given a --url and --depth the application returns Links to the specified depth.

It crawls the given '--url' for any links and keeps on discovering new links within the found links till the '--depth'

For example: if the max depth is 2, and the main url is `http://example.com/` and there is a chain of links from `/` -> `/foo` -> `/bar` -> `/baz` -> `/quux`
depth 0: crawl `/` and discover the link to `/foo`
depth 1: crawl `/foo` and discover `/bar`
depth 2: crawl `/bar` and discover `/baz`
stops at depth 2, does not crawl `/baz` and therefore will not discover the link to `/quux`
 Collect the links with their base url to produce a list:
1. http://example.com/
2. http://example.com/foo
3. http://example.com/bar
4. http://example.com/baz

