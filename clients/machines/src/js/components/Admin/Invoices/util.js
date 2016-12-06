function statusInfo(invoice) {
  var info = '';

  if (invoice.get('FastbillNo') || invoice.get('Status') === 'outgoing') {
    info += 'Invoice';
  } else {
    info += 'Draft';
  }

  if (invoice.get('Canceled')) {
    if (invoice.get('CanceledSent')) {
      info += ' (Canceled & Cancelation Sent)';
    } else {
      info += ' (Canceled & Cancelation Unsent)';
    }
  } else {
    if (invoice.get('Sent')) {
      info += ' (Sent)';
    } else {
      info += ' (Unsent)';
    }
  }

  return info;
}

export default {
  statusInfo
};
