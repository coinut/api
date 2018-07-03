import hmac
import hashlib
import json
import uuid
import random
import time
import requests


class CoinutAPI():
    '''REST API for https://coinut.com. More documents can be find at https://github.com/coinut/api/wiki'''

    def __init__(self, user = None, api_key = None):
        '''Initialize the API

        Args:
            user (str): your username
            api_key (str): your REST API Key on https://coinut.com/account/settings
        '''
        self.user = user
        self.api_key = api_key

    def get_balance(self):
        '''Get my balance

        Returns:
            your balance in a dict.

        Examples:
            >>> c = CoinutAPI('your username', 'your REST API Key')
            >>> print c.get_balance()
            {u'USDT': u'18775.86604171', u'status': [u'OK'], u'XMR': u'0.00154154', u'BTC': u'63.51457338', ...}

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-account-balance
        '''
        return self.request("user_balance")



    def get_spot_instruments(self, pair = None):
        '''Get spot trading instruments information

        Args:
            pair (str): it can be any spot trading pair like "BTCUSDT"
            or "LTCBTC".

        Returns:
            if pair argument is specified, return the pair's
            information in a dict; otherwise returns all spot trading
            pairs' information

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments

        '''

        result = self.request("inst_list", {'sec_type': 'SPOT'})
        if pair != None:
            return result['SPOT'][pair][0]
        else:
            return result['SPOT']



    def get_spot_inst_id(self, pair):
        '''Get a spot trading instrument's inst_id. This inst_id is needed for
        submitting, canceling, or querying orders and for checking ticks and
        orderbooks.

        Args:
            pair (str): it can be any spot trading pair like "BTCUSDT" or "LTCBTC".

        Returns:
            return the spot trading pair's inst_id.

        Examples:
            >>> c = CoinutAPI()
            >>> print c.get_spot_inst_id('LTCBTC')
            1

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments

        '''

        return self.get_spot_instruments(pair)['inst_id']



    def get_inst_tick(self, inst_id):
        '''Get a spot trading instrument's last tick.

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

        Returns:
            the instrument's last tick

        Examples:
            >>> c = CoinutAPI()
            >>> print c.get_inst_tick(1)
            {u'nonce': 289676251, u'low24': u'0.01235000', u'last': u'0.01253000', ... u'trans_id': 4312849734}

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-realtime-ticks
        '''

        return self.request("inst_tick", {"inst_id": inst_id})


    def get_orderbook(self, inst_id):
        '''Get a spot trading instrument's orderbook.

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

        Returns:
            the trading pair's orderbook

        Examples:
            >>> c = CoinutAPI()
            >>> print c.get_orderbook(1)
            {u'nonce': 1605656906, u'sell': [{u'count': 1, u'price': u'0.01311', u'qty': u'0.00200000'...}

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-orderbooks-in-realtime
        '''

        return self.request("inst_order_book", {"inst_id": inst_id})



    def get_open_orders(self, inst_id):
        '''Get my open orders.

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

        Returns:
            my open orders for an instrument

        Examples:
            >>> c = CoinutAPI('your username', 'your REST API Key')
            >>> print c.get_open_orders(1)
            [{u'order_id': 1120194747, u'open_qty': u'37.08640000', u'price': u'0.01251000', u'qty': u'37.08640000'...

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#get-open-orders
        '''

        return self.request("user_open_orders", {"inst_id": inst_id})['orders']


    def create_new_order(self, inst_id, side, qty, price = None, client_ord_id = None):
        '''Create a dict containing the information for opening a new order

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

            side (str): either 'BUY' or 'SELL'. It's case sensitive.

            qty (float): the quantity that you want to buy or sell

            price (float): if price is None, it's a market order; otherwise it's the price of the limit order.

            client_ord_id (int): an id specified by the client.

        Returns:
            a dict containing the information for the new order

        Examples:
            >>> c = CoinutAPI()
            >>> print c.create_new_order(1, 'BUY', 0.0000001, 0.013)
            {'price': '0.01300000', 'qty': '0.00000010', 'side': 'BUY', 'client_ord_id': 1170372055, 'inst_id': 1}
        '''

        order = {'inst_id': inst_id, "side": side, 'qty': "%.8f" % qty}
        if price is not None:
            order['price'] = "%.8f" % price

        if client_ord_id is not None:
            order['client_ord_id'] = client_ord_id
        else:
            order['client_ord_id'] = random.randint(1, 4294967290)
        return order


    def submit_new_order(self, inst_id, side, qty, price = None, client_ord_id = None):
        '''Submit an order to the exchange

        Args:
            inst_id (int): the inst_id can be obtained using the get_spot_inst_id or get_spot_instruments functions.

            side (str): either 'BUY' or 'SELL'. It's case sensitive.

            qty (float): the quantity that you want to buy or sell

            price (float): use None to specifcy that the order is a market order; otherwise it's the price of the limit order.

            client_ord_id (int): an id specified by the client.

        Returns:
            An order can be rejected, partially filled, or filled. Check https://github.com/coinut/api/wiki/Websocket-API#submit-an-order for the formats

        Examples:
            >>> c = CoinutAPI()
            >>> print c.submit_new_order(1, 'BUY', 0.0000001, 0.012, 598847915)

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#submit-an-order
        '''
        return self.request("new_order", self.create_new_order(inst_id, side, qty, price, client_ord_id))


    def submit_new_orders(self, ords):
        '''Submit a batch of orders to the exchange

        Args:
            ords (list): a list of orders. Every order may be
            generated using the create_new_order function. The maximum
            number of orders is 1000.

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#submit-an-order
        '''
        return self.request("new_orders", {"orders": ords})


    def cancel_order(self, inst_id, order_id):
        '''Cancel an order

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

            order_id (int): the id of the order to be canceled

        Returns:
            Order cancel result

        Examples:
            >>> c = CoinutAPI('your username', 'your REST API Key')
            >>> print c.cancel_order(1, 3355)

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#cancel-an-order

        '''
        return self.request("cancel_order", {'inst_id': inst_id, 'order_id': order_id})


    def cancel_orders(self, inst_id, order_ids):
        '''Cancel orders in a batch

        Args:
            inst_id (int): the inst_id can be obtained using the
            get_spot_inst_id or get_spot_instruments functions.

            order_ids (a list of int): all the ids of the orders to be
            canceled. The maximum number of orders is 1000.

        Returns:
            cancel results

        Examples:
            >>> c = CoinutAPI('your username', 'your REST API Key')
            >>> print c.cancel_orders(1, [3355, 1345])

        See also:
            https://github.com/coinut/api/wiki/Websocket-API#cancel-orders-in-batch-mode
        '''

        ords = [{'inst_id': inst_id, 'order_id': x} for x in order_ids]
        return self.request("cancel_orders", {'entries': ords})


    def request(self, api, content = {}):
        url = 'https://api.coinut.com'
        content["request"] = api
        content["nonce"] = random.randint(1, 4294967290)
        content = json.dumps(content)
        headers = {}
        if self.api_key is not None and self.user is not None:
            sig = hmac.new(self.api_key, msg=content,
                           digestmod=hashlib.sha256).hexdigest()
            headers = {'X-USER': self.user, "X-SIGNATURE": sig}

        response = requests.post(url, headers=headers, data=content, timeout=5)
        return response.json()
