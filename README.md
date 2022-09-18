# Crypto 💸🌐 Web Scraper & Exporter

This project is utilized to scrape a web endpoints to utilize client data rather than utilize a paid API. Freemium at the cost of my own hardware computing power and storage.


## Features

- Select URLs From LiveCoinWatch.com for coin metrics
- Prometheus Exporter
- Dockerized

## Deployment

To deploy this project run

```bash
  git clone <project.git>
  sudo docker-compose up -d
```
## Add More Metrics

Append the url section with additional urls using the YAML array string syntax shown below with a hypen and the url in quotes. 
Please note you will need to target the url of a specific page with a coin based on how the scraper is designed.

```
    metrics:
    crypto:
        url: 
        - "https://www.livecoinwatch.com/price/Avalanche-AVAX"
        - "https://www.livecoinwatch.com/price/Ethereum-ETH"
        - "https://www.livecoinwatch.com/price/Bitcoin-BTC"
        type: gauge
        description: "Exports the current coin value based on the web scrape"

```

## Authors

- [@Aetrius](https://www.github.com/Aetrius)


## Roadmap

- Additional Query Metrics

## Related Sources
[Live Coin Watch](https://www.livecoinwatch.com/price/Cardano-ADA)

[Web Scrape Resource](https://www.scrapingbee.com/blog/web-scraping-go/)

[Colly Web Scrape Docs](http://go-colly.org/)
