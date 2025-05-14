local config = import 'default.jsonnet';

config {
  'mrmintchain_9000-1'+: {
    'app-config'+: {
      'minimum-gas-prices': '100000000000mnt',
    },
    genesis+: {
      app_state+: {
        feemarket+: {
          params+: {
            base_fee:: super.base_fee,
          },
        },
      },
    },
  },
}
