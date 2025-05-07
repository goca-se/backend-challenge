from pydantic import BaseModel, Field
from typing import List, Optional
from enum import Enum


class TipoServico(str, Enum):
    TODOS = "todos"
    NORMAL = "normal"
    EXPRESSO = "expresso"


class Remetente(BaseModel):
    cep: str


class Destinatario(BaseModel):
    cep: str


class Pacote(BaseModel):
    peso: float
    altura: float
    largura: float
    comprimento: float
    valor_declarado: float


class Servico(BaseModel):
    tipo: TipoServico = TipoServico.TODOS


class CotacaoRequest(BaseModel):
    remetente: Remetente
    destinatario: Destinatario
    pacote: Pacote
    servico: Servico


class OpcaoResponse(BaseModel):
    tipo: str
    codigo: str
    valor: float
    prazo: int
    codigo_rastreio_modelo: str


class OpcoesResponse(BaseModel):
    opcao: List[OpcaoResponse]


class CotacaoSuccessResponse(BaseModel):
    status: str = "sucesso"
    opcoes: OpcoesResponse
    mensagem: str = "Cotação realizada com sucesso"


class CotacaoErrorResponse(BaseModel):
    status: str = "erro"
    codigo_erro: str
    mensagem: str 