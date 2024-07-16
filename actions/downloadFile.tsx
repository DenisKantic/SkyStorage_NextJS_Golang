"use server";

export async function downloadFile(id: number) {
    const response = await fetch(`http://localhost:8080/download?id=${id}`, {
        method: 'GET',
        mode: 'cors'
    });

    if (!response.ok) {
        throw new Error('Failed to fetch file');
    }

    const blob = await response.blob();
    return { blob, contentDisposition: response.headers.get('Content-Disposition') };
}
