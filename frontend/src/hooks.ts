import { useEffect } from "react"

export function usePageTitle(pageTitle?: string) {
    useEffect(() => {
        if (pageTitle) {
            document.title = `${pageTitle} - Tiny URL`;
        } else {
            document.title = "Tiny URL";
        }
    }, [pageTitle]);
}
