document.getElementById("signupForm").addEventListener("submit", function(event) {
    event.preventDefault();
    var formData = new FormData(this);
    fetch("/users", {
        method: "POST",
        body: formData
    })
    .then(response => {
    
        if (!response.ok) {
            throw new Error("Erro ao cadastrar usuário.");
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            window.location.href = "/outra-pagina";
        } else {
            alert("Erro: " + data.error);
        }
    })
    .catch(error => {
        console.error("Erro ao enviar requisição:", error);
        alert("Erro ao cadastrar usuário.");
    });
});