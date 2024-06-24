import React from 'react';
import { useLoaderData } from 'react-router-dom';
import { DateTime } from 'luxon';
import { sessionGET, errorFromResponse } from '../api.js';

import SAMUELHangar from '../components/page/SAMUELHangar.js';

export async function dashboardLoader() {
    let loaderData = {};

    const apiResponse = await sessionGET('dashboard');
    if (apiResponse.ok) {
        const {user, session, [user.role.id]: userDetails} = await apiResponse.json();
        loaderData = {
            user: {
                firstName: userDetails.firstName,
                role: user.role
            },
            session: {
                expiresOn: DateTime.fromISO(session.expiresOn)
            },
        };
    } else {
        switch (apiResponse.status) {
        case 401:
            throw await errorFromResponse(apiResponse, 'Your Session Has Expired', {path: '/login', name: 'Login Page'})
        default:
            throw await errorFromResponse(apiResponse, 'Failed To Load Dashboard', {path: '/login', name: 'Login Page'});
        }
    }

    return loaderData;
} 

export default function DashboardPage() {
    const loaderData = useLoaderData();
    
    return (
        <SAMUELHangar user={loaderData.user} session={loaderData.session}>
            <h3>Welcome, {loaderData.user.firstName}!</h3>
        </SAMUELHangar>
    );
}