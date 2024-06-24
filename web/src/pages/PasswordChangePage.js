import React from 'react';
import { Link, useActionData, redirect, useParams } from 'react-router-dom';
import { object, ref, string } from 'yup';
import { POST, PUT, errorFromResponse } from '../api.js';

import Button from 'react-bootstrap/Button';
import ListGroup from 'react-bootstrap/ListGroup';

import SAMUELAlert from '../components/control/SAMUELAlert.js';
import SAMUELForm from '../components/form/SAMUELForm.js';
import SAMUELLander from '../components/page/SAMUELLander.js';

export async function passwordChangeAction({ params, request: actionRequest }) {
    let actionData = {};

    const {token} = params;
    const formData = Object.fromEntries(await actionRequest.formData());
    if (!token) {
        const apiResponse = await POST('password_change/create', {email: formData.email});
        if (apiResponse.ok) {
            switch (apiResponse.status){
            case 201:
                actionData.alert = {
                    variant: 'success',
                    body: 'Please check your inbox to complete your password change request.'
                };
                break;
            case 200:
                actionData.alert = {
                    variant: 'warning',
                    body: 'Your password change request is still pending. Please check your inbox.'
                };
                break;
            }
        } else {
            throw await errorFromResponse(apiResponse, 'Failed to Change Password', {path: '/login', name: 'Login Page'});
        }
    } else {
        const apiResponse = await PUT(`password_change/fulfill/${token}`, {newPassword: formData.newPassword});
        if (apiResponse.ok) {
            actionData.alert = {
                variant: 'success',
                body: 'You have succesfully changed your password.',
                redirect: {
                    path: '/login',
                    name: 'Login Page'
                }
            };
        } else {
            throw await errorFromResponse(apiResponse, 'Failed to Request Password Change', { path: '/login', name: 'Login Page' });
        }
    }

    return actionData;
}

export default function PasswordChangeRequestPage() {
    const {token} = useParams();
    
    const {alert} = useActionData() ?? {alert: null};

    if (!token) {
        const schema = object()
            .shape({
                email: string()
                    .required('Email is required.')
                    .email('Email is invalid.')
            });

        return (
            <SAMUELLander title="Need Some Help?" subtitle="Request A Password Change">
                <ListGroup variant="flush">
                    <ListGroup.Item>
                        <h3 className="my-3">For Administrators, Instructors, and Students:</h3>
                        <p className="my-3">
                            Contact the IT Department at <Link to="mailto:myit@southhills.edu">myit@southhills.edu</Link>.
                        </p>
                    </ListGroup.Item>
                    <ListGroup.Item>
                        <h3 className="my-3">For Supervisors:</h3>
                        <p>Please enter your email below:</p>
                        <SAMUELForm action="/password_change" validationSchema={schema}>
                            <SAMUELForm.Control type="email" name="email" label="Email" />
                            <Button type="submit" variant="primary" size="lg">Request Password Change</Button>
                            <SAMUELAlert variant={alert?.variant} show={!!alert}>{alert?.body}</SAMUELAlert>
                        </SAMUELForm>
                    </ListGroup.Item>
                </ListGroup>
            </SAMUELLander>
        );
    } else {
        const schema = object()
            .shape({
                newPassword: string()
                    .required('New password is required.'),
                newPasswordConfirmation: string()
                    .required('New password confirmation is required')
                    .equals([ref('newPassword')], 'Passwords do not match.')
            });

        return (
            <SAMUELLander title="Welcome to SAMUEL!" subtitle="Change Your Password">
                <SAMUELForm action={`/password_change/${token}`} validationSchema={schema}>
                    <SAMUELForm.Control type="password" name="newPassword" label="New Password" />
                    <SAMUELForm.Control type="password" name="newPasswordConfirmation" label="Confirm New Password" />
                    <Button type="submit" variant="primary" size="lg">Change Password</Button>
                    <SAMUELAlert variant={alert?.variant} show={!!alert} redirect={alert?.redirect}>{alert?.body}</SAMUELAlert>
                </SAMUELForm>
            </SAMUELLander>
        );
    }
}