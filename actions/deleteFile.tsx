"use server"

export async function deleteFile(id: number){
    const response = await fetch(`http://localhost:8080/delete?id=${id}`,{
        method: "DELETE"
    })

    const data = await response.json();

    return data;
}