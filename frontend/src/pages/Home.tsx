import axios from "axios";
import { FormEvent, useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardHeader, CardContent, CardFooter } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

function TinyUrl() {
    const [tinyUrl, setTinyUrl] = useState<string>();
    const [loading, setLoading] = useState<boolean>();
    const [error, setError] = useState<string>();

    const validateURL = useCallback((input: string) => {
        // Check if the entered URL is valid and has a valid
        // http(s) scheme. Otherwise, the backend redirects
        // will not work properly.
        let url;
        try {
            url = new URL(input);
        } catch {
            return false
        }
        return url.protocol === "http:" || url.protocol === "https:";
    }, [])

    const handleSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setError(undefined);
        const formData = new FormData(event.currentTarget);
        if (!validateURL(formData.get("url") as string)) {
            setError("Please enter a valid URL starting with https.");
            return;
        }

        // make backend request and handle response
        setLoading(true);
        axios.postForm('/api/shortenurl', formData)
            .then(function (response) {
                setLoading(false);
                const shortUrl = response.data["tinyurl"];
                setTinyUrl(shortUrl);
            })
            .catch(function (error) {
                setError(error as string);
                setTinyUrl(undefined);
                setLoading(false);

            })
    }, []);

    const copyToClipboard = useCallback(() => {
        navigator.clipboard.writeText(tinyUrl ?? "");
    }, [tinyUrl]);

    const navigateToShorURL = useCallback(() => {
        window.open(tinyUrl, "_blank");
    }, [tinyUrl]);

    const renderResult = () => {
        return <div className="border-l ml-4">
            <div className="ml-4 space-y-4">
                {
                    loading ?
                    <div className="space-y-2">
                        <Skeleton className="h-8 w-[250px]" />
                        <Skeleton className="h-8 w-[250px]" />
                        <Skeleton className="h-8 w-[250px]" />
                    </div> :
                    tinyUrl ?
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="tinyurl">Short URL</Label>
                            <Input
                                readOnly
                                id="tinyurl"
                                type="text"
                                value={tinyUrl}
                            />
                        </div>
                        <div className="grid grid-cols-2 gap-2 pt-4">
                            <Button title="Copy to Clipboard" onClick={copyToClipboard}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <rect width="8" height="4" x="8" y="2" rx="1" ry="1"/>
                                    <path d="M8 4H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2"/>
                                    <path d="M16 4h2a2 2 0 0 1 2 2v4"/>
                                    <path d="M21 14H11"/>
                                    <path d="m15 10-4 4 4 4"/>
                                </svg>
                                Copy to clipboard
                            </Button>
                            <Button title="Go to link" onClick={navigateToShorURL}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M15 3h6v6"/>
                                    <path d="M10 14 21 3"/>
                                    <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                                </svg>
                                Go to link
                            </Button>
                        </div>
                    </div> : null
                }
            </div>
        </div>
    };

    return (
        <Card>
            <CardHeader className="text-center">
                <h3 className="text-lg font-semibold">Create a tiny URL</h3>
            </CardHeader>
            <CardContent>
                <div className={tinyUrl ? "grid grid-cols-2" : ""}>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="longurl">Long URL</Label>
                            <Input id="longurl" name="url" />
                        </div>
                        {error ? <p className="text-sm">{error}</p> : null}
                        <div className="pt-4">
                            <Button className="w-full" type="submit">Create</Button>
                        </div>
                    </form>
                    {
                        renderResult()
                    }
                </div>
            </CardContent>
            <CardFooter />
        </Card>
    );
}

export function Home() {
    return (
        <div className="flex flex-col gap-4">
            <TinyUrl />
        </div>
    );
}