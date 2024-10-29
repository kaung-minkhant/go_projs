# Content

This project is a sitemap crawler, given an initial sitemap. 
It crawls through every available sitemaps, and crawls the pages.
It then get some SEO related data from each of the pages.

This project heavily rely on **channels** and **go routines** to make simultaneous request
and process the results from all of the **go routines**. It uses **goquery** package to
parse response data, such as HTML parsing.
