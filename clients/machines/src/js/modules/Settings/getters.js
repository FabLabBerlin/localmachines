const getAdminSettings = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.get('settings');
  }
];

const getCurrency = [
  ['settingsStore'],
  (settingsStore) => {
    return settingsStore.getIn(['settings', 'Currency', 'ValueString']);
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
    return settingsStore.getIn(['settings', 'VAT', 'ValueFloat']);
  }
];

export default {
	getAdminSettings,
  getCurrency,
  getFastbillTemplates,
  getVatPercent
};
