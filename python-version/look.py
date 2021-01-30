from datetime import datetime, timedelta
import json
import os
import requests
import time
from bs4 import BeautifulSoup


MAX_CHECKING_HOURS = 5
MIN_CHECK_BACKOFF_SECONDS = 60

LOOP_TIMEOUT_HOURS = 5
LOOP_RETRY_SECONDS = 60


class Operator:
    _config_path = ''
    _duration = 5
    _backoff = 60
    _items = []
    _created_at = datetime.now()

    def __init__(self, path):
        self._config_path = path
        self._created_at = datetime.now()
        data = {}
        with open(path, 'r') as f:
            data = json.load(f)

        items = []
        for i in data["zara_items"]:
            items.append(ZaraObject(i['url'], i['sizes']))

        self._set(data["checking_time_hours"], data["check_every_seconds"], items)

    def _set(self, duration, backoff, items):
        self._duration = duration
        if self._duration > MAX_CHECKING_HOURS:
            self._duration = MAX_CHECKING_HOURS
        self._backoff = backoff
        if self._backoff < MIN_CHECK_BACKOFF_SECONDS:
            self._backoff = MIN_CHECK_BACKOFF_SECONDS
        self._items = items

    def get_duration(self):
        return self._duration

    def get_backoff(self):
        return self._backoff

    def get_items(self):
        return self._items

    def update_items(self):
        data = {}
        with open(self._config_path, 'r') as f:
            data = json.load(f)
        self._items = parse_zara_items(data["zara_items"])


def parse_zara_items(arr):
    items = []
    for i in arr:
        items.append(ZaraObject(i['url'], i['sizes']))
    return items


class ZaraObject:
    url = ""
    sizes = []
    name = ""

    def __init__(self, url, sizes):
        self.url = url
        self.sizes = sizes

    def __str__(self):
        return "'{}' in sizes: '{}'".format(self.name, self.sizes)

    def init_name(self, name):
        self.name = name

    def check_availability(self, sizes):
        for desired in self.sizes:
            val = sizes.get(desired, False)
            if val:
                message = "{}".format(self.name, desired, self.url)
                notify2("ZARA SIZE AVAILABLE", "Size {}".format(desired), message, self.url)

    def search(self):
        page = requests.get(self.url)
        soup = BeautifulSoup(page.text, 'html.parser')

        # get product name
        self.name = soup.find('h1', 'product-name').text.rstrip("Dane szczegółowe")

        print("checking", self.name)

        # <label class="product-size _product-size" data-name="XL" data-sku="78476118">
        # or
        # <label class="product-size _product-size disabled _disabled" data-name="XXL" data-sku="78476117">
        sizes_labels = soup.find('div', 'size-list').find_all('label', 'product-size')
        sizes = {}
        for s in sizes_labels:
            sizes[s['data-name']] = 'disabled' not in s['class']

        self.check_availability(sizes)


def notify(title, text):
    os.system("""
        osascript -e 'say "Hey Zara has your size!"'
        osascript -e 'display dialog "{}" message "{}"'
    """.format(title, text))


def notify2(title, subtitle, message, url):
    os.system("""
        osascript -e 'say "Checkout ZARA"'
        terminal-notifier -title '{}' -subtitle '{}' -message '{}' -open '{}'
    """.format(title, subtitle, message, url))


if __name__ == '__main__':
    op = Operator('./config.json')

    started_at = datetime.now()
    while True:
        for i in op.get_items():
            i.search()

        if datetime.now() > started_at + timedelta(hours=op.get_duration()):
            print("finishing checking")
            break
        else:
            print("sleeping for {}".format(op.get_backoff()))
            time.sleep(op.get_backoff())
