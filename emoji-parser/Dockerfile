# Starting from the Openjdk-8 container
FROM java:openjdk-8-jdk

# Set the WORKDIR. All following commands will be run in this directory.
WORKDIR /app

# Copying all gradle files necessary to install gradle with gradlew
COPY gradle /app/gradle
COPY build.gradle /app/build.gradle
COPY settings.gradle /app/settings.gradle
COPY gradlew /app/gradlew

ADD src /app/src

RUN ./gradlew fatJar

# Install the gradle version used in the repository through gradlew


CMD ["java","-jar","build/libs/emoji-all-1.0-SNAPSHOT.jar"]