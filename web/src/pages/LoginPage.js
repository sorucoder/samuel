import React from 'react';
import { Link, redirect, useLoaderData, useActionData } from 'react-router-dom';
import { object, string } from 'yup';
import { credentialsGET, errorFromResponse, errorResponse } from '../api.js';

import Button from 'react-bootstrap/Button';

import SAMUELAlert from '../components/control/SAMUELAlert.js';
import SAMUELForm from '../components/form/SAMUELForm.js';
import SAMUELLander from '../components/page/SAMUELLander.js';

export function loginLoader({request}) {
    const url = new URL(request.url);
    if (url.searchParams.get('expired')) {
        return {
            alert: {
                variant: 'danger',
                body: 'Your session has expired. Please login.'
            }
        }
    }

    return null;
}

export async function loginAction({request}) {
    const formData = Object.fromEntries(await request.formData());
    const apiResponse = await credentialsGET('login', formData.identity, formData.password);
    if (!apiResponse.ok) {
        switch (apiResponse.status) {
        case 401:
            return {
                alert: {
                    variant: 'danger',
                    body: 'Your credentials are invalid. Please try again.'
                }
            }
        default:
            throw await errorFromResponse(apiResponse, 'Failed To Login', {path: '/login', name: 'Login Page'})
        }
    } else {
        const {session} = await apiResponse.json();
        localStorage.setItem('samuel_session_token', session.token);
        return redirect(`/dashboard`);
    }
}

export default function LoginPage() {
    const schema = object()
        .shape({
            identity: string()
                .required('Username or email is required.')
                .matches(
                    /^(?:[a-z]+[0-9]*|[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$/,
                    { message: 'Username or email is invalid.' }
                ),
            password: string()
                .required('Password is required.')
        });
    
    const {alert} = useLoaderData() ?? useActionData() ?? {alert: null};

    return (
        <SAMUELLander title="Welcome to S.A.M.U.E.L.!" subtitle="Login">
            <SAMUELForm action="/login" validationSchema={schema}>
                <SAMUELForm.Control type="text" name="identity" label="Username or Email" />
                <SAMUELForm.Control type="password" name="password" label="Password" />
                <Link to="/password_change">Forgot Password?</Link>
                <Button type="submit" variant="primary" size="lg">Login</Button>
                <SAMUELAlert variant={alert?.variant} show={!!alert}>{alert?.body}</SAMUELAlert>
            </SAMUELForm>
        </SAMUELLander>
    );
};