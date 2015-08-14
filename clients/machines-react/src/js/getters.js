const getIsLogged = [
	['loginStore'],
	(loginStore) => {
		return loginStore.get('isLogged');
	}
];

const getUid = [
	['loginStore'],
	(loginStore) => {
		return loginStore.get('uid');
	}
];

const getFirstTry = [
	['loginStore'],
	(loginStore) => {
		return loginStore.get('firstTry');
	}
];

export default { getIsLogged, getUid, getFirstTry };
