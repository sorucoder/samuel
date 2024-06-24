import React, { useState } from 'react';

import FloatingLabel from 'react-bootstrap/FloatingLabel';
import FormControl from 'react-bootstrap/FormControl';
import Feedback from 'react-bootstrap/Feedback';

export default function SAMUELControl({ type, name, label, value, schema, onChange }) {
    const [isInvalid, setIsInvalid] = useState(false);
    const onBlur = async () => {
        try {
            await schema.validate();
        } catch (errors) {
            setIsInvalid(true);
            return;
        }
        setIsInvalid(false);
    };

    return (
        <FloatingLabel controlId={`${name}InteractiveFormControl`} label={label}>
            <FormControl type={type} name={name} value={value} isInvalid={isInvalid} placeholder={label} onBlur={onBlur} onChange={onChange} />
            <Feedback type="invalid"></Feedback>
        </FloatingLabel>
    );
}