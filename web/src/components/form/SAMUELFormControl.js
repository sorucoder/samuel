import React from 'react';
import { useFormikContext } from 'formik';

import FloatingLabel from 'react-bootstrap/FloatingLabel';
import FormControl from 'react-bootstrap/FormControl';
import Feedback from 'react-bootstrap/Feedback';

export default function SAMUELFormControl({type, name, label}) {
    const formik = useFormikContext();

    return (
        <FloatingLabel controlId={`${name}FormControl`} label={label}>
            <FormControl type={type} name={name} isValid={formik.touched[name] && !formik.errors[name]} isInvalid={formik.touched[name] && formik.errors[name]} placeholder={label} onBlur={formik.handleBlur} onChange={formik.handleChange} />
            <Feedback type={formik.errors[name] ? 'invalid' : 'valid'}>{formik.errors[name]}</Feedback>
        </FloatingLabel>
    );
}