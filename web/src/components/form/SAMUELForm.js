import React from 'react';
import { useSubmit } from 'react-router-dom';
import { Formik } from 'formik';

import Form from 'react-bootstrap/Form';

import SAMUELFormControl from './SAMUELFormControl.js';

function SAMUELForm({validationSchema, initialValues = {}, action, children}) {
    initialValues ??= {};

    const submit = useSubmit();

    return (
        <Formik validationSchema={validationSchema} initialValues={initialValues} onSubmit={(values) => submit(values, {method: 'POST', action})}>
            {({validated, handleSubmit}) => (
                <Form className="d-grid gap-2" noValidate validated={validated} onSubmit={handleSubmit}>
                    {children}
                </Form>
            )}
        </Formik>
    );
}

SAMUELForm.Control = SAMUELFormControl;

export default SAMUELForm;