# Build stage
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Copy csproj and restore
COPY ["Klokr_Dashboard.csproj", "./"]
RUN dotnet restore "Klokr_Dashboard.csproj"

# Copy everything and build
COPY . .
RUN dotnet build "Klokr_Dashboard.csproj" -c Release -o /app/build

# Publish stage
FROM build AS publish
RUN dotnet publish "Klokr_Dashboard.csproj" -c Release -o /app/publish /p:UseAppHost=false

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app
EXPOSE 8080

COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Klokr_Dashboard.dll"]

