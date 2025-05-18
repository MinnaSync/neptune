import { okAsync, errAsync } from "neverthrow";

export async function request(
    url: string,
    method: 'get' | 'post' | 'put' | 'delete',
    headers: Record<string, string>,
    query: Record<string, string>,
) {
    const reqUrl = new URL(url);

    for (const [key, value] of Object.entries(query)) {
        reqUrl.searchParams.set(key, value);
    }

    const req = await fetch(reqUrl, { method, headers });

    if (req.status !== 200) {
        return errAsync(new Error('Failed to request.'));
    }

    const contentType = req.headers.get('content-type');
    if (contentType?.includes('application/json')) {
        const json = await req.json();
        return okAsync(json);
    }

    if (contentType?.startsWith('text')) {
        const text = await req.text();
        return okAsync(text);
    }

    return okAsync(null); // dunno what to do with this for now
}