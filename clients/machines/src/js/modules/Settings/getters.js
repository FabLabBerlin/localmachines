const getAdminSettings = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('adminSettings');
  }
];

const getVatPercent = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('VatPercent');
  }
];

export default {
	getAdminSettings,
  getVatPercent
};
