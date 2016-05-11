const getVatPercent = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('VatPercent');
  }
];

export default {
  getVatPercent
};
