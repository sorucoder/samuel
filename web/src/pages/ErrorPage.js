import React from 'react';
import { useRouteError } from 'react-router-dom';

import Image from 'react-bootstrap/Image';

import SAMUELAlert from '../components/control/SAMUELAlert';
import SAMUELLander from '../components/page/SAMUELLander.js';

export default function ErrorPage() {
    const error = useRouteError();

    if (error.status) {
        if (error.details) {
            return (
                <SAMUELLander title="Keep Calm, Interns!" subtitle="Something Went Wrong">
                    <Image className="w-50 align-self-center" src={`https://http.dog/${error.status}.jpg`} />
                    <hr />
                    <SAMUELAlert variant="danger" show>
                        {error.message}:<br />
                        <code>{error.details}</code>
                    </SAMUELAlert>
                </SAMUELLander>
            );
        } else {
            return (
                <SAMUELLander title="Keep Calm, Interns!" subtitle="Something Went Wrong">
                    <Image className="w-50 align-self-center" src={`https://http.dog/${error.status}.jpg`} />
                    <hr />
                    <SAMUELAlert variant="danger" show redirect={error.redirect}>
                        {error.message}
                    </SAMUELAlert>
                </SAMUELLander>
            );
        }
    } else {
        return (
            <SAMUELLander title="Keep Calm, Interns!" subtitle="Something Went Wrong">
                <SAMUELAlert variant="danger" show redirect={error.redirect}>
                    {error.message}
                </SAMUELAlert>
            </SAMUELLander>
        );
    }
}