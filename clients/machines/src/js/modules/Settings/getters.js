const getAdminSettings = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('adminSettings');
  }
];

const getFastbillTemplates = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('fastbillTemplates');
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
  getFastbillTemplates,
  getVatPercent
};
