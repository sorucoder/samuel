/**
 * Makes an unauthorized GET request to the API.
 *
 * @export
 * @param {string} path - The path of the API method, relative to /api.
 * @param {null|object} parameters - The query parameters to be passed to the GET request, if any.
 * @returns {Promise<Response>}
 */
export function GET(path, parameters = null) {
    if (parameters) {
        const query = new URLSearchParams(parameters);
        return fetch(`/api/${path}?${query}`, {method: 'GET'});
    } else {
        return fetch(`/api/${path}`, { method: 'GET' });
    }
}

/**
 * Makes an unauthorized POST request to the API.
 *
 * @export
 * @param {string} path - The path of the API method, relative to /api.
 * @param {any} payload - The payload of the POST request.
 * @returns {Promise<Response>}
 */
export function POST(path, payload) {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const body = JSON.stringify(payload);
    return fetch(`/api/${path}`, {method: 'POST', headers, body});
}

/**
 * Makes an unauthorized PUT request to the API.
 *
 * @export
 * @param {string} path - The path of the API method, relative to /api.
 * @param {any} payload - The body of the POST request.
 * @returns {Promise<Response>}
 */
export function PUT(path, payload) {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const body = JSON.stringify(payload);
    return fetch(`/api/${path}`, {method: 'PUT', headers, body});
}


/**
 * Makes an authorized GET request to the API using the provided credentials.
 *
 * @export
 * @param {string} path- The path of the API method, relative to /api.
 * @param {string} identity - The username or email of the user.
 * @param {string} password - The password of the user.
 * @returns {Promise<Response>}
 */
export function credentialsGET(path, identity, password) {
    const headers = new Headers();
    const credentials = btoa(`${identity}:${password}`);
    headers.append('Authorization', `Basic ${credentials}`);
    return fetch(`/api/${path}`, {method: 'GET', headers});
}

/**
 * Makes an authorized GET request to the API using the current session token found in local storage.
 * If there is no session token found in local storage, this function throws an error which, if uncaught,
 * triggers the error boundary and will redirect them back to /login.
 *
 * @export
 * @param {string} path - The path of the API method, relative to /api.
 * @param {null|object} parameters - The query parameters to be passed to the GET request, if any.
 * @throws An error boundary object that will redirect to /login.
 * @returns {Promise<Response>}
 */
export function sessionGET(path, parameters = null) {
    const sessionToken = localStorage.getItem('samuel_session_token');
    if (!sessionToken) {
        throw {
            message: 'Your Session Has Expired',
            redirect: {
                path: '/login',
                name: 'Login Page'
            }
        };
    }
    const headers = new Headers();
    headers.append('Authorization', `Bearer ${sessionToken}`);
    if (parameters) {
        const query = new URLSearchParams(parameters);
        return fetch(`/api/${path}?${query}`, {method: 'GET', headers});
    } else {
        return fetch(`/api/${path}`, {method: 'GET', headers});
    }
}

/**
 * Generates an error object from an API response.
 * This function will pass the status from the API response to the new error object.
 * If the API response did reach the API server and return a JSON response, its "details" field will be passed.
 * Otherwise, it will assume it is a text response and set details to that.
 *
 * @export
 * @async
 * @param {Promise<Response>} response - The API response to generate the error from.
 * @param {string} message - The user-friendly error message.
 * @param {null|object} redirect - An optional redirect 
 * @returns {object}
 */
export async function errorFromResponse(response, message, redirect = null) {
    const {headers} = response;
    if (headers.get('Content-Type').includes('application/json')) {
        const body = await response.json();
        if (body.details) {
            return {
                status: response.status,
                message: message,
                details: `${body.error}: ${body.details}`,
                redirect: redirect
            };
        } else {
            return {
                status: response.status,
                message: message,
                redirect: redirect
            };
        }
    } else {
        const body = await response.text();
        return {
            status: response.status,
            message: message,
            details: body,
            redirect: redirect
        };
    }
}