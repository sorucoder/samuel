import React from 'react';
import { useLoaderData, redirect } from "react-router-dom";
import { sessionGET, errorFromResponse } from '../api';

import SAMUELAlert from '../components/control/SAMUELAlert.js';
import SAMUELLander from '../components/page/SAMUELLander.js';

export async function logoutLoader() {
    let loaderData = {};
    
    const apiResponse = await sessionGET('logout');
    if (apiResponse.ok) {
        loaderData.alert = {
            variant: 'success',
            body: 'You have logged out.',
            redirect: {
                path: '/login',
                name: 'Login Page'
            }
        };
    } else {
        switch (apiResponse.status) {
        case 401:
            return redirect('/login');
        default:
            throw await errorFromResponse(apiResponse, 'Failed To Logout', {path: '/dashboard', name: 'Dashboard'})
        }
    }

    return loaderData;
}

export default function LogoutPage() {
    const {alert} = useLoaderData() ?? {alert: null};

    return (
        <SAMUELLander title="Thank you for using S.A.M.U.E.L.!" subtitle="Logout">
            <SAMUELAlert variant={alert?.variant} show redirect={alert?.redirect}>{alert?.body}</SAMUELAlert>
        </SAMUELLander>
    );
}