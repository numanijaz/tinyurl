import axios from "axios";
import { FormEvent, useCallback, useState } from "react"

export function UrlShortner() {
    const [shortUrl, setShortUrl] = useState();
    const [error, setError] = useState();
    const onFormSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const params = {
            "url": formData.get("longurl")
        }
        axios.post('/shortenurl', params)
            .then(function (response) {
                const shortUrl = response.data["tinyurl"];
                setShortUrl(shortUrl);
            })
            .catch(function (error) {
                setError(error);
            });
    }, []);

    return (
        <div className="flex flex-row justify-center items-center">
            <div>
                <form id="shortenUrlForm" className="w-full mt-8 space-y-6">
                    <div id="shortenUrlError" className="hidden bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded"></div>
                    <div className="rounded-md shadow-sm -space-y-px">
                        <div>
                            <label htmlFor="longurl" className="block text-sm font-medium text-gray-700 mb-2 px-1">Long URL</label>
                            <input 
                                id="longurl" 
                                name="longurl" 
                                type="text"
                                required
                                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm" 
                                placeholder="Your long URL"
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            className="relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        >
                            <span id="shortenUrl">Shorten URL</span>
                        </button>
                    </div>
                </form>
            </div>
            <div>
                
            </div>
        </div>
    )
}