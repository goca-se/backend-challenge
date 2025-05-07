from fastapi import FastAPI, Request, Depends, HTTPException, Response
from fastapi.responses import PlainTextResponse
from app.auth import verify_credentials
from app.models import (
    CotacaoRequest, 
    CotacaoSuccessResponse, 
    CotacaoErrorResponse, 
    OpcaoResponse, 
    OpcoesResponse
)
import asyncio
import random
import lxml.etree as ET
from typing import Union
import time

app = FastAPI(
    title="EntregasRápidas API Mock",
    description="API mock para simulação da transportadora EntregasRápidas",
    version="1.0.0"
)

# Global variable to track valid CEPs for demonstration purposes
# VALID_CEPS = [
#     "01310100", "04538132", "05422010", "22041080", "30130170", 
#     "40010000", "50030220", "60170001", "70200002", "80010010"
# ]

def validate_cep(cep: str) -> bool:
    return int(cep) >= 50000000 and int(cep) <= 99999999


async def parse_xml_request(request: Request) -> CotacaoRequest:
    """
    Parse the XML request body into a CotacaoRequest model.
    """
    try:
        body = await request.body()
        root = ET.fromstring(body)
        
        remetente_cep = root.xpath('//remetente/cep')[0].text
        destinatario_cep = root.xpath('//destinatario/cep')[0].text
        
        pacote_peso = float(root.xpath('//pacote/peso')[0].text)
        pacote_altura = float(root.xpath('//pacote/altura')[0].text)
        pacote_largura = float(root.xpath('//pacote/largura')[0].text)
        pacote_comprimento = float(root.xpath('//pacote/comprimento')[0].text)
        pacote_valor_declarado = float(root.xpath('//pacote/valor_declarado')[0].text)
        
        servico_tipo = root.xpath('//servico/tipo')[0].text
        
        return CotacaoRequest(
            remetente={"cep": remetente_cep},
            destinatario={"cep": destinatario_cep},
            pacote={
                "peso": pacote_peso,
                "altura": pacote_altura,
                "largura": pacote_largura,
                "comprimento": pacote_comprimento,
                "valor_declarado": pacote_valor_declarado
            },
            servico={"tipo": servico_tipo}
        )
    except Exception as e:
        raise HTTPException(
            status_code=400,
            detail=f"Erro ao processar XML: {str(e)}"
        )


def model_to_xml(
    model: Union[CotacaoSuccessResponse, CotacaoErrorResponse]
) -> str:
    """
    Convert a Pydantic model to XML string.
    """
    root = ET.Element("resultado")
    
    status = ET.SubElement(root, "status")
    status.text = model.status
    
    if model.status == "sucesso":
        # Success response
        success_model = model  # type: CotacaoSuccessResponse
        
        opcoes = ET.SubElement(root, "opcoes")
        for opcao_data in success_model.opcoes.opcao:
            opcao = ET.SubElement(opcoes, "opcao")
            
            tipo = ET.SubElement(opcao, "tipo")
            tipo.text = opcao_data.tipo
            
            codigo = ET.SubElement(opcao, "codigo")
            codigo.text = opcao_data.codigo
            
            valor = ET.SubElement(opcao, "valor")
            valor.text = f"{opcao_data.valor:.2f}"
            
            prazo = ET.SubElement(opcao, "prazo")
            prazo.text = str(opcao_data.prazo)
            
            codigo_rastreio = ET.SubElement(opcao, "codigo_rastreio_modelo")
            codigo_rastreio.text = opcao_data.codigo_rastreio_modelo
        
        mensagem = ET.SubElement(root, "mensagem")
        mensagem.text = success_model.mensagem
    else:
        # Error response
        error_model = model  # type: CotacaoErrorResponse
        
        codigo_erro = ET.SubElement(root, "codigo_erro")
        codigo_erro.text = error_model.codigo_erro
        
        mensagem = ET.SubElement(root, "mensagem")
        mensagem.text = error_model.mensagem
    
    # Generate XML with declaration and pretty printing
    xml_declaration = '<?xml version="1.0" encoding="UTF-8"?>\n'
    xml_string = ET.tostring(root, pretty_print=True, encoding="unicode")
    
    return xml_declaration + xml_string


@app.post("/api/v1/entregas-rapidas/cotacao", response_class=PlainTextResponse)
async def cotacao_frete(
    request: Request,
    username: str = Depends(verify_credentials)
):
    """
    Endpoint para cotação de fretes da transportadora EntregasRápidas.
    
    Recebe solicitação em formato XML e retorna as opções disponíveis.
    
    Simula latência de 2 a 5 segundos.
    """
    # Parse XML request
    cotacao_request = await parse_xml_request(request)
    
    # Simulate high latency (2-5 seconds)
    delay = random.uniform(2.0, 5.0)
    await asyncio.sleep(delay)
    
    # Check if CEP is valid
    if (
        not validate_cep(cotacao_request.remetente.cep) or 
        not validate_cep(cotacao_request.destinatario.cep)
    ):
        error_response = CotacaoErrorResponse(
            codigo_erro="ER-005",
            mensagem="CEP de destino inválido ou não atendido"
        )
        return Response(
            content=model_to_xml(error_response),
            media_type="application/xml"
        )
    
    # Prepare success response
    options = []
    
    # Add normal option
    options.append(
        OpcaoResponse(
            tipo="normal",
            codigo="ER-NORMAL",
            valor=random.uniform(20.0, 40.0),  # Random price between 20-40
            prazo=random.randint(2, 4),  # Random delivery time between 2-4 days
            codigo_rastreio_modelo="ER1234567890BR"
        )
    )
    
    # Add express option
    options.append(
        OpcaoResponse(
            tipo="expresso",
            codigo="ER-EXPRESS",
            valor=random.uniform(40.0, 60.0),  # Random price between 40-60
            prazo=1,  # Express is always 1 day
            codigo_rastreio_modelo="EX1234567890BR"
        )
    )
    
    success_response = CotacaoSuccessResponse(
        opcoes=OpcoesResponse(opcao=options)
    )
    
    return Response(
        content=model_to_xml(success_response),
        media_type="application/xml"
    )


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "service": "entregas-rapidas-api"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app.main:app", host="0.0.0.0", port=3000, reload=True) 