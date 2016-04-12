const getVatPercent = [
  ['settingsStore'],
  (settingsStore) => {
    console.log('settingsStore:', settingsStore);
    return settingsStore.get('VatPercent');
  }
];

export default {
  getVatPercent
};
