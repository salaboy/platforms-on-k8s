'use client'
import { useState, useEffect } from 'react'
import styles from '@/app/styles/events.module.css'
import Textfield from '../forms/textfield/textfield';
import Button from '../forms/button/button';
import Select from '../forms/select/select';
import Switch from '../forms/switch/switch';


function NewEnvironment() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [sended, setSended] = useState(false);
    const [name, setName] = useState("");
    const [type, setType] = useState("");
    const [installInfra, setInstallInfra] = useState(false);
    const [debug, setDebug] = useState(false);

    const handleBack = () => {
        setSended(false)
      }

    const handleSubmit = () => {
        setLoading(true);
        setIsError(false);

        var name = document.getElementById("name").value
        console.log("Name: " + name)
        var e = document.getElementById("type");
        var envType = e.value;
        console.log("Type: " + envType)
        var installInfra = document.getElementById("installInfra").checked
        console.log("Install Infra? " + installInfra)

        var debug = document.getElementById("debug").checked
        console.log("debug? " + debug)

        const data = {
            name: name,
            parameters: {
                type: envType,
                installInfra: installInfra,
                frontend: {
                    debug: debug,
                },
            }
        }
        
        console.log("Sending Post!" + JSON.stringify(data))
        try {
            fetch('/api/environments/', {
                method: "POST",
                body: JSON.stringify(data),
                headers: {
                    'accept': 'application/json',
                },
            }).then((response) => response.json())
                .then((data) => {
                    setName('');
                    setType('');
                    setInstallInfra('');
                    setDebug('');
                    setLoading(false);
                    setSended(true);
                })
        } catch (err) {
            setLoading(false);
            setIsError(true);
        }
    }

    

    return (
        <div>


            <div className={styles.EventsList}>
                {!sended && (
                    <div>

                        <Textfield label="Name" id="name" name="name" />
                        <Select label="Type" id="type" name="type">
                            <option value="development">Development</option>
                        </Select>
                        <Switch label="Install Infrastructure" id="installInfra" name="installInfra" />
                        <Switch label="Frontend Debug" id="debug" name="debug" />
                        

                        {isError && <small className="mt-3 d-inline-block text-danger">Something went wrong. Please try again later.</small>}

        
                        <Button type="submit" clickHandler={handleSubmit} >Create New Environment</Button>
        
                    </div>
                )}
                {sended && (
                    <>
                        <h3>Thanks for creating an new Environment with us!</h3>
                        <Button clickHandler={handleBack} >Create another Environment</Button>
                    </>
                )}
            </div>
        </div>
        
    );

}
export default NewEnvironment;