function promisify(fn) {
	// return a new promisified function
	return (...args) => {
		return new Promise((resolve, reject) => {
			// create a callback that resolves and rejects
			function callback(err, result) {
				if (err) {
					reject(err);
				} else {
					resolve(result);
				}
			}

			args.push(callback);

			// pass the callback into the function
			fn.call(this, ...args);
		});
	};
}

function jwtDecode(token) {
	return JSON.parse(atob(token.split(".")[1]))
}

export { promisify, jwtDecode };
