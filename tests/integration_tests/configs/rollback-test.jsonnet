local config = import 'default.jsonnet';

config {
  'mrmintchain_9000-1'+: {
    validators: super.validators[0:1] + [{
      name: 'fullnode',
    }],
  },
}
