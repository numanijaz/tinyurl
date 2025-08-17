import axios from "axios";
import { useCallback, useEffect } from "react";
import { useParams } from "react-router-dom";

export function NavigateTinyUrl() {
    const {hash} = useParams();
    useEffect(() => {
        axios.get(`/${hash}`)
            .then(data => {
                console.log(data);
            })
            .catch(error => {
                console.log(error);
            });
    }, [hash]);

    return "";
}
